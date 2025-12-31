package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSubscription_IsString(t *testing.T){
	dataSt := "test"

	res := isString(dataSt)

	require.True(t, res)
}

func TestSubscription_IsUUID(t *testing.T){
	dataUUID := uuid.New()

	res := isUUID(dataUUID)

	require.True(t, res)
}

func TestSubscription_IsUint32(t *testing.T){
	dataUint32 := uint32(1)

	res := isUint32(dataUint32)

	require.True(t, res)
}

func TestSubscription_IsTime(t *testing.T){
	dataTime := time.Now()

	res := isTime(dataTime)

	require.True(t, res)
}
