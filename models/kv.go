package models

type KV struct {
	Id    int64
	Key   string `xorm:"unique index"`
	Value []byte
}

func (x *Engine) Set(key string, value []byte) error {
	kv := new(KV)
	has, e := x.Where("key = ?", key).Get(kv)
	if !has {
		kv.Key = key
		kv.Value = value
		e = x.Save(kv)
		if e != nil {
			return e
		}
	} else {
		kv.Value = value
		_, e = x.Id(kv.Id).Update(kv)
		if e != nil {
			return e
		}
	}
	return nil
}

func (x *Engine) Get(key string) ([]byte, error) {
	kv := new(KV)
	has, e := x.Where("key = ?", key).Get(kv)
	if !has {
		return nil, e
	} else {
		return kv.Value, nil
	}
	return nil, nil
}
