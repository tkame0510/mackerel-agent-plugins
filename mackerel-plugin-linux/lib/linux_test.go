package mplinux

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectWho(t *testing.T) {
	_, err := os.Stat("/usr/bin/who")
	if err != nil {
		return
	}
	p := make(map[string]interface{})

	assert.Nil(t, collectWho(&p))
}

func TestParseWho(t *testing.T) {
	stub := `test0  pts/48       2014-09-30 08:00 (192.168.24.123)
test1  pts/48       2014-09-30 08:59 (192.168.24.123)
test2  pts/48       2014-09-30 09:00 (192.168.24.123)`
	stat := make(map[string]interface{})

	err := parseWho(stub, &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["users"], 3)
}

func TestParseWho2(t *testing.T) {
	stub := ""
	stat := make(map[string]interface{})

	err := parseWho(stub, &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["users"], 0)
}

func TestGetWho(t *testing.T) {
	_, err := os.Stat("/usr/sbin/who")
	if err != nil {
		return
	}

	ret, err := getWho()
	assert.Nil(t, err)
	assert.NotNil(t, ret)
}

func TestCollectStat(t *testing.T) {
	path := "/proc/stat"
	_, err := os.Stat(path)
	if err != nil {
		return
	}
	p := make(map[string]interface{})

	assert.Nil(t, collectProcStat(path, &p))
}

func TestParseProcStat(t *testing.T) {
	stub := `intr 614818624 122 8 0 0 1 0 0 0 1 0 0 0 123 0 0 0 0 0 0 0 0 0 0 0 4846888 0 44650320 253 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
 ctxt 879305394
 btime 1409212617
 processes 1959410`
	stat := make(map[string]interface{})

	err := parseProcStat(bytes.NewBufferString(stub), &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["interrupts"], 614818624)
	assert.EqualValues(t, stat["context_switches"], 879305394)
	assert.EqualValues(t, stat["forks"], 1959410)
}

func TestCollectNetworkStat(t *testing.T) {
	_, err := os.Stat("/usr/sbin/ss")
	if err != nil {
		return
	}
	p := make(map[string]interface{})

	assert.Nil(t, collectNetworkStat(&p))
}

func TestParseSs(t *testing.T) {
	stub := `State      Recv-Q Send-Q                       Local Address:Port                         Peer Address:Port 
LISTEN     0      128                                     :::45103                                  :::*     
LISTEN     0      128                                     :::111                                    :::* 
TIME-WAIT  0      0                         ::ffff:127.0.0.1:80                       ::ffff:127.0.0.1:50082 
ESTAB      0      0                              10.0.25.101:60826                         10.0.25.104:5672  `
	stat := make(map[string]interface{})

	err := parseSs(bytes.NewBufferString(stub), &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["LISTEN"], 2)
	assert.EqualValues(t, stat["TIME-WAIT"], 1)
	assert.EqualValues(t, stat["ESTAB"], 1)
}

func TestParseSs2(t *testing.T) {
	stub := `Netid State      Recv-Q Send-Q                                      Local Address:Port                                        Peer Address:Port
nl    UNCONN     0      0                                                      18:0                                                       *
p_raw UNCONN     0      0                                                       *:em2                                                     *
u_dgr UNCONN     0      0                                                /dev/log 10549                                                  * 0
u_dgr LISTEN     0      0                                       /run/udev/control 8552                                                   * 0
u_str LISTEN     0      10                                  /var/run/acpid.socket 9649                                                   * 0
u_str ESTAB      0      0                                    @/com/ubuntu/upstart 10582                                                  * 1887`
	stat := make(map[string]interface{})

	err := parseSs(bytes.NewBufferString(stub), &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["LISTEN"], 2)
	assert.EqualValues(t, stat["UNCONN"], 3)
	assert.EqualValues(t, stat["ESTAB"], 1)
}

func TestCollectProcVmstat(t *testing.T) {
	path := "/proc/vmstat"
	_, err := os.Stat(path)
	if err != nil {
		return
	}
	p := make(map[string]interface{})

	assert.Nil(t, collectProcVmstat(path, &p))
}

func TestParseProcVmstat(t *testing.T) {
	stub := `pgpgin 770294
pgpgout 31351354
pswpin 0
pswpout 113`
	stat := make(map[string]interface{})

	err := parseProcVmstat(bytes.NewBufferString(stub), &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["pgpgin"], 770294)
	assert.EqualValues(t, stat["pgpgout"], 31351354)
	assert.EqualValues(t, stat["pswpin"], 0)
	assert.EqualValues(t, stat["pswpout"], 113)
}

func TestCollectDiskStats(t *testing.T) {
	path := "/sys"

	_, err := os.Stat(path)
	if err != nil {
		return
	}
	p := make(map[string]interface{})

	assert.Nil(t, collectDiskStats(path, &p))
}

func TestParseDiskStat(t *testing.T) {
	name := "testdevice"
	stub := `  36049      277  3702446    36470  1165021   131631 15197712  1648460        0   771090  1684180`
	stat := make(map[string]interface{})

	err := parseDiskStat(name, stub, &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat[fmt.Sprintf("iotime_%s", name)], 771090)
	assert.EqualValues(t, stat[fmt.Sprintf("iotime_weighted_%s", name)], 1684180)
	assert.EqualValues(t, stat[fmt.Sprintf("tsreading_%s", name)], 36470)
	assert.EqualValues(t, stat[fmt.Sprintf("tswriting_%s", name)], 1648460)
}
