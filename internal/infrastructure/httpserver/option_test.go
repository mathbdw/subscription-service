package httpserver

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOption_Address(t *testing.T){
	h := "localhost"
	p := uint16(1111)
	
	s := &Server{}
	opt := Address(h, p)
	opt(s)

	require.Equal(t, fmt.Sprintf("%s:%d", h, p), s.address)
}

func TestOption_Prefork(t *testing.T){
	s := &Server{}
	opt := Prefork(true)
	opt(s)

	require.Equal(t, true, s.prefork)
}

func TestOption_ReadTimeout(t *testing.T){
	tm := 5 * time.Second
	s := &Server{}
	opt := ReadTimeout(tm)
	opt(s)

	require.Equal(t, tm, s.readTimeout)
}

func TestOption_WriteTimeout(t *testing.T){
	tm := 5 * time.Second
	s := &Server{}
	opt := WriteTimeout(tm)
	opt(s)

	require.Equal(t, tm, s.writeTimeout)
}

func TestOption_ShutdownTimeout(t *testing.T){
	tm := 5 * time.Second
	s := &Server{}
	opt := ShutdownTimeout(tm)
	opt(s)

	require.Equal(t, tm, s.shutdownTimeout)
}