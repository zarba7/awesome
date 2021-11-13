package base

type Option struct {
	NameSpace string
	Sid int32
	Release  bool
	Version  string
	RocketMQ struct{
		NameServer []string
	}
	Application struct {
		Name string
		RPCAddr string
		Param   param
	}
}

type param map[string]interface{}

func (the param) Find(keys ...string) param {
	var val interface{}
	var ok bool
	val = the
	for _, k := range keys {
		val, ok = val.(param)[k]
		if !ok {
			return nil
		}
	}
	return val.(map[string]interface{})
}

func (the param) find(k string, keys []string) interface{} {
	if len(keys) == 0 {
		return the[k]
	}
	val := the.Find(k).Find(keys[:len(keys)-1]...)
	k = keys[len(keys)-1]
	return val[k]
}

func (the param) Int(k string, keys ...string) int { return int(the.Int64(k, keys...)) }
func (the param) Str(k string, keys ...string) string {
	val := the.find(k, keys)
	if val == nil {
		return ""
	}
	return val.(string)
}
func (the param) Int64(k string, keys ...string) int64 { return the.find(k, keys).(int64) }
func (the param) Bool(k string, keys ...string) bool {
	val := the.find(k, keys)
	return val != nil && val.(bool)
}