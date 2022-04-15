package utils

import (
	"regexp"
	"strconv"

	"github.com/marmotedu/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var sizeMatcher = regexp.MustCompile(`.*\((\d+)\s*vs\.\s*(\d+)\)`)

type MsgAssembler func(msg []byte) interface{}

func StreamSendMsg(stream grpc.ServerStream, data []byte, chunksize int, assembler MsgAssembler) error {
	i := 0
	for i < len(data) {
		var msg []byte
		if i+chunksize < len(data) {
			msg = data[i : i+chunksize]
		} else {
			msg = data[i:]
		}
		// TODO: Waiting for the improvement of the inspection mechanism of gRPC server max send msg size
		//s1, s2, err := resolveResExhaustedErr(stream.SendMsg(assembler(msg)))
		//if err != nil {
		//	return err
		//}
		//if s1-s2 > 0 {
		//	chunksize -= s1 - s2
		//	continue
		//}
		err := stream.SendMsg(assembler(msg))
		if err != nil {
			return err
		}
		i += chunksize
	}
	return nil
}

func resolveResExhaustedErr(err error) (int, int, error) {
	if status.Code(err) == codes.ResourceExhausted {
		if m := sizeMatcher.FindStringSubmatch(err.Error()); len(m) == 3 {
			s1, e1 := strconv.Atoi(m[1])
			s2, e2 := strconv.Atoi(m[2])
			err = errors.NewAggregate([]error{e1, e2})
			if err != nil {
				return 0, 0, errors.Wrap(err, "failed to resolve valid message size")
			}
			return s1, s2, nil
		} else {
			return 0, 0, errors.Wrap(err, "failed to resolve gRPC `ResourceExhausted` error")
		}
	}
	return 0, 0, err
}
