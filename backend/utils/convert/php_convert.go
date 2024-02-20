package convutil

import (
	"os/exec"
)

type PhpConvert struct {
	CmdConvert
}

const phpDecodeCode = `
<?php

$action = strtolower($argv[1]);
$content = $argv[2];

if ($action === 'decode') {
    $decoded = base64_decode($content);
    if ($decoded !== false) {
        $obj = unserialize($decoded);
        if ($obj !== false) {
            $unserialized = json_encode($obj, JSON_UNESCAPED_UNICODE);
            if ($unserialized !== false) {
                echo base64_encode($unserialized);
                return;
            }
        }
    }
} elseif ($action === 'encode') {
    $decoded = base64_decode($content);
    if ($decoded !== false) {
        $json = json_decode($decoded, true);
        if ($json !== false) {
            $serialized = serialize($json);
            if ($serialized !== false) {
                echo base64_encode($serialized);
                return;
            }
        }
    }
}
echo '[RDM-ERROR]';
`

func NewPhpConvert() *PhpConvert {
	c := CmdConvert{
		Name:       "PHP",
		Auto:       true,
		DecodePath: "php",
		EncodePath: "php",
	}

	var err error
	if err = exec.Command(c.DecodePath, "-v").Err; err != nil {
		return nil
	}

	var filepath string
	if filepath, err = c.writeExecuteFile([]byte(phpDecodeCode), "php_decoder.php"); err != nil {
		return nil
	}
	c.DecodeArgs = []string{filepath, "decode"}
	c.EncodeArgs = []string{filepath, "encode"}

	return &PhpConvert{
		CmdConvert: c,
	}
}

func (p *PhpConvert) Enable() bool {
	if p == nil {
		return false
	}
	return true
}

func (p *PhpConvert) Encode(str string) (string, bool) {
	if !p.Enable() {
		return str, false
	}
	return p.CmdConvert.Encode(str)
}

func (p *PhpConvert) Decode(str string) (string, bool) {
	if !p.Enable() {
		return str, false
	}
	return p.CmdConvert.Decode(str)
}
