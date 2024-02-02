package server

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"time"

	"word_of_wisdom/internal/pkg/data"
	"word_of_wisdom/internal/pkg/pow"

	"github.com/google/uuid"
)

var (
	ErrExit = fmt.Errorf("client close connection")

	wordOfWisdom = []string{
		"If any of you lacks wisdom, let him ask of God, who gives to all liberally and without reproach, and it will be given to him",
		"For I will give you words and a wisdom that none of your opponents will be able to withstand or contradict",
		"For to one is given the word of wisdom through the Spirit, to another the word of knowledge through the same Spirit",
		"But they could not stand up against his wisdom or the Spirit by whom he spoke",
		"Conduct yourselves wisely toward outsiders, making the most of the time",
		"If any of you lacks wisdom, he should ask God, who gives generously to all without finding fault, and it will be given to him",
		"Who is wise and understanding among you? Let him show by good conduct that his works are done in the meekness of wisdom",
		"But the wisdom from above is first pure, then peaceable, gentle, open to reason, full of mercy and good fruits, impartial and sincere",
	}
)

type cache interface {
	Add(uuid.UUID)
	Get(uuid.UUID) bool
	Delete(uuid.UUID)
}

type Server struct {
	cache           cache
	listener        net.Listener
	firstZerosCount int
}

func NewServer(cache cache, address string, firstZero int) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("server.NewServer: %w", err)
	}

	return &Server{
		listener:        listener,
		cache:           cache,
		firstZerosCount: firstZero,
	}, nil
}

func (s *Server) Start() error {
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("server.Start: %w", err)
		}

		slog.Info(fmt.Sprintf("connection from %s", conn.RemoteAddr()))

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			slog.Error("server.handleConnection - reader: %v", err)
			return
		}
		data, err := s.processRequest(req, conn.RemoteAddr().String())
		if err != nil {
			slog.Error("server.handleConnection - request: %v", err)
			return
		}
		if data != nil {
			err := sendData(*data, conn)
			if err != nil {
				slog.Warn("server.handleConnection: %v", err)
			}
		}
	}
}

func (s *Server) processRequest(dataStr string, clientInfo string) (*data.Data, error) {
	d, err := data.Parse(dataStr)
	if err != nil {
		return nil, fmt.Errorf("server.ProcessRequest: %w", err)
	}
	switch d.Key {
	case data.Challenge:
		randValue := uuid.New()
		s.cache.Add(randValue)

		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: s.firstZerosCount,
			Date:       time.Now().Unix(),
			Resource:   clientInfo,
			Rand:       base64.RawStdEncoding.EncodeToString([]byte(randValue.String())),
			Counter:    0,
		}
		hashcashMarshaled, err := json.Marshal(hashcash)
		if err != nil {
			return nil, fmt.Errorf("server.ProcessRequest - marshal hashcash: %v", err)
		}
		d := data.Data{
			Key:   data.Challenge,
			Value: string(hashcashMarshaled),
		}
		return &d, nil
	case data.Response:
		var hashcash pow.HashcashData
		err := json.Unmarshal([]byte(d.Value), &hashcash)
		if err != nil {
			return nil, fmt.Errorf("server.ProcessRequest - unmarshal hashcash: %w", err)
		}
		if hashcash.Resource != clientInfo {
			return nil, fmt.Errorf("server.ProcessRequest - invalid hashcash resource")
		}

		randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
		if err != nil {
			return nil, fmt.Errorf("server.ProcessRequest - decode: %w", err)
		}
		randValue, err := uuid.ParseBytes(randValueBytes)
		if err != nil {
			return nil, fmt.Errorf("server.ProcessRequest - parse: %w", err)
		}

		exists := s.cache.Get(randValue)
		if !exists {
			return nil, fmt.Errorf("server.ProcessRequest - not found in cache")
		}

		maxIter := hashcash.Counter
		if maxIter <= 0 {
			maxIter = 1
		}
		_, err = hashcash.ComputeHashcash(maxIter)
		if err != nil {
			return nil, fmt.Errorf("server.ProcessRequest - invalid hashcash")
		}
		d := data.Data{
			Key:   data.Grant,
			Value: wordOfWisdom[rand.Intn(len(wordOfWisdom))],
		}
		s.cache.Delete(randValue)
		return &d, nil
	default:
		return nil, fmt.Errorf("server.ProcessRequest - unknown key: %s", d.Key)
	}
}

func sendData(d data.Data, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", d.String())
	_, err := conn.Write([]byte(msgStr))
	return err
}
