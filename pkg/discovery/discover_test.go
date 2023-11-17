package discovery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscover_PingENode(t *testing.T) {
	enode, err := parseNode("enode://ccc5ae45b01c3cbdcdd15275af187144416718a208d66d41d2ecf5cd874b13f834a022dd589163bb759f7a8907ee7987ba5a29f003347f57178d720a2998575b@13.215.34.27:5050")
	assert.Nil(t, err)
	fmt.Printf("Node info: %+v \n", enode)

}

func TestDiscover_CrawlRound(t *testing.T) {

}
