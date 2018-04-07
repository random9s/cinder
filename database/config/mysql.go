package config

import (
	"encoding/json"
)

//MySQL ...
// Check https://godoc.org/github.com/go-sql-driver/mysql#Config to see all options
type MySQL struct {
	User                    string            `json:"user"`
	Passwd                  string            `json:"password"`
	Net                     string            `json:"net"`
	Addr                    string            `json:"address"`
	DBName                  string            `json:"name"`
	Params                  map[string]string `json:"params"`
	Collation               string            `json:"collation"`
	Loc                     string            `json:"time_location"`
	MaxAllowedPacket        int               `json:"max_allowed_packet"`
	TLSConfig               string            `json:"tls_config"`
	AllowAllFiles           bool              `json:"allow_all_files"`
	AllowCleartextPasswords bool              `json:"allow_cleartext_passwords"`
	AllowNativePasswords    bool              `json:"allow_native_passwords"`
	AllowOldPasswords       bool              `json:"allow_old_passwords"`
	ClientFoundRows         bool              `json:"client_found_rows"`
	ColumnsWithAlias        bool              `json:"columns_with_alias"`
	InterpolateParams       bool              `json:"interpolate_params"`
	MultiStatements         bool              `json:"multi_statements"`
	ParseTime               bool              `json:"parse_time"`
	RejectReadOnly          bool              `json:"reject_read_only"`
	Strict                  bool              `json:"strict"`
}

//Register returns the mysql driver config file
func (m *MySQL) Register(data []byte) error {
	return json.Unmarshal(data, &m)
}
