package local

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/marmotedu/errors"
)

func PosBasedOnConfig(ctx context.Context) (*v1.ContextPos, error) {
	pos := context.GetPos(ctx)
	father := pos.Target().Father()
	_, childIndex := pos.Position()
	if err := father.Error(); err != nil {
		return nil, err
	}

	if father.Type() == context_type.TypeConfig {
		config, ok := father.(*Config)
		if !ok {
			return nil, errors.WithCode(code.ErrV3InvalidContext, "father context is not a config")
		}

		return &v1.ContextPos{
			ConfigPath: config.FullPath(),
			PosIndex:   []int32{int32(childIndex)},
		}, nil
	} else if father.Type() == context_type.TypeMain {
		main, ok := father.(*Main)
		if !ok {
			return nil, errors.WithCode(code.ErrV3InvalidContext, "father context is not a main config")
		}

		return &v1.ContextPos{
			ConfigPath: main.MainConfig().FullPath(),
			PosIndex:   []int32{int32(childIndex)},
		}, nil
	}

	ctxPos, err := PosBasedOnConfig(father)
	if err != nil {
		return nil, err
	}
	ctxPos.PosIndex = append(ctxPos.PosIndex, int32(childIndex))

	return ctxPos, nil
}
