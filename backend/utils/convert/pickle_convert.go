package convutil

type PickleConvert struct {
	CmdConvert
}

const pickleDecodeCode = `
import base64
import json
import pickle
import sys

if __name__ == "__main__":
    if len(sys.argv) >= 3:
        action = sys.argv[1].lower()
        content = sys.argv[2]

        try:
            if action == 'decode':
                decoded = base64.b64decode(content)
                obj = pickle.loads(decoded)
                unserialized = json.dumps(obj, ensure_ascii=False)
                print(base64.b64encode(unserialized.encode('utf-8')).decode('utf-8'))
            elif action == 'encode':
                decoded = base64.b64decode(content)
                obj = json.loads(decoded)
                serialized = pickle.dumps(obj)
                print(base64.b64encode(serialized).decode('utf-8'))
        except:
            print('[RDM-ERROR]')
    else:
        print('[RDM-ERROR]')

`

func NewPickleConvert() *PickleConvert {
	c := CmdConvert{
		Name: "Pickle",
		Auto: true,
	}
	c.DecodePath, c.EncodePath = "python3", "python3"
	var err error
	if _, err = runCommand(c.DecodePath, "--version"); err != nil {
		c.DecodePath, c.EncodePath = "python", "python"
		if _, err = runCommand(c.DecodePath, "--version"); err != nil {
			return nil
		}
	}
	// check if pickle available
	if _, err = runCommand(c.DecodePath, "-c", "import pickle"); err != nil {
		return nil
	}
	var filepath string
	if filepath, err = writeExecuteFile([]byte(pickleDecodeCode), "pickle_decoder.py"); err != nil {
		return nil
	}
	c.DecodeArgs = []string{filepath, "decode"}
	c.EncodeArgs = []string{filepath, "encode"}

	return &PickleConvert{
		CmdConvert: c,
	}
}

func (p *PickleConvert) Enable() bool {
	if p == nil {
		return false
	}
	return true
}

func (p *PickleConvert) Encode(str string) (string, bool) {
	if !p.Enable() {
		return str, false
	}
	return p.CmdConvert.Encode(str)
}

func (p *PickleConvert) Decode(str string) (string, bool) {
	if !p.Enable() {
		return str, false
	}
	return p.CmdConvert.Decode(str)
}
