package main

type Config struct {
	NetworkID uint64
}

func U2UMainnet() Config {
	return Config{
		NetworkID: uint64(1234),
	}
}
