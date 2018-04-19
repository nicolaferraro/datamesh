package protobuf

func NewUpsertOperation(path string, version uint64, data []byte) *Operation {
	return &Operation{
		&Operation_Upsert{
			&UpsertOperation{
				&Data{
					Path:    &Path{path, version},
					Content: data,
				},
			},
		},
	}
}

func NewDeleteOperation(path string, version uint64) *Operation {
	return &Operation{
		&Operation_Delete{
			&DeleteOperation{
				&Path{path, version},
			},
		},
	}
}

func NewReadOperation(path string, version uint64) *Operation {
	return &Operation{
		&Operation_Read{
			&ReadOperation{
				&Path{path, version},
			},
		},
	}
}
