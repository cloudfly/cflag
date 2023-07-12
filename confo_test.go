package confo

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type testConfig struct {
	APPName string `default:"confgo" json:",omitempty"`
	Hosts   []string

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306" json:",omitempty"`
		SSL      bool   `default:"true" json:",omitempty"`
	}

	Contact struct {
		Name  string
		Email string `required:"true"`
	} `env:"-"`

	Boolean bool   `arg:"b,bool"`
	Number  int64  `arg:"n,num"`
	Help    string `arg:"h,help"`

	private string
}

func generateFile(t *testing.T, c any) string {
	file, err := os.CreateTemp("/tmp", "confo")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
		return ""
	}
	defer file.Close()
	content, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
		return ""
	}
	if _, err := file.Write(content); err != nil {
		t.Fatal(err)
		t.FailNow()
		return ""
	}
	return file.Name()
}

func generateDefaultConfig() testConfig {
	return testConfig{
		APPName: "confgo",
		Hosts:   []string{"http://example.org", "http://hello.world"},
		DB: struct {
			Name     string
			User     string `default:"root"`
			Password string `required:"true" env:"DBPassword"`
			Port     uint   `default:"3306" json:",omitempty"`
			SSL      bool   `default:"true" json:",omitempty"`
		}{
			Name:     "confgo",
			User:     "confgo",
			Password: "confgo",
			Port:     3306,
			SSL:      true,
		},
		Contact: struct {
			Name  string
			Email string `required:"true"`
		}{
			Name:  "Cloudfly",
			Email: "hello@gmail.com",
		},
		Number: 100,
	}
}

func TestLoadNormaltestConfig(t *testing.T) {
	config := generateDefaultConfig()
	f := generateFile(t, config)
	defer os.Remove(f)
	var result testConfig
	err := Load(&result, f)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(result, config) {
		t.Errorf("result should equal to original configuration")
	}
}

func TestDefaultValue(t *testing.T) {
	config := generateDefaultConfig()
	config.APPName = ""
	config.DB.Port = 0
	config.DB.SSL = false
	f := generateFile(t, config)
	defer os.Remove(f)

	var result testConfig
	err := Load(&result, f)

	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(result, generateDefaultConfig()) {
		t.Errorf("result should be set default value correctly")
	}
}

func TestMissingRequiredValue(t *testing.T) {
	config := generateDefaultConfig()
	config.DB.Password = ""

	f := generateFile(t, config)
	defer os.Remove(f)

	var result testConfig
	err := Load(&result, f)

	t.Log(result.DB.Password)

	if err == nil {
		t.Errorf("result should be return error because the password requried")
	}
}

func TestLoadtestConfigurationByEnvironment(t *testing.T) {
	config := generateDefaultConfig()
	config2 := struct {
		APPName string
	}{
		APPName: "config2",
	}
	f1 := generateFile(t, config)
	defer os.Remove(f1)
	f2 := generateFile(t, config2)
	os.Rename(f2, f1+".test")
	f2 = f1 + ".test"
	defer os.Remove(f2)

	var result testConfig
	if err := New(&Config{Env: "test"}).Load(&result, f1); err != nil {
		t.Errorf("No error should happen when load configurations, but got %v", err)
	}

	var defaultConfig = generateDefaultConfig()
	defaultConfig.APPName = "config2"
	if !reflect.DeepEqual(result, defaultConfig) {
		t.Errorf("result should be load configurations by environment correctly")
	}
}

func TestOverwritetestConfigurationWithEnvironmentWithDefaultPrefix(t *testing.T) {
	config := generateDefaultConfig()
	f := generateFile(t, config)
	defer os.Remove(f)

	os.Setenv("CONFGO_APPNAME", "config2")
	os.Setenv("CONFGO_HOSTS", "- http://example.org\n- http://hello.world")
	os.Setenv("CONFGO_DB_NAME", "db_name")
	os.Setenv("CONFGO_NAME", "cloudfly22")
	defer os.Setenv("CONFGO_APPNAME", "")
	defer os.Setenv("CONFGO_HOSTS", "")
	defer os.Setenv("CONFGO_DB_NAME", "")
	defer os.Setenv("CONFGO_NAME", "")

	config.APPName = "config2"
	config.Hosts = []string{"http://example.org", "http://hello.world"}
	config.DB.Name = "db_name"
	config.Contact.Name = "cloudfly22"

	var result testConfig
	if err := New(&Config{EnvPrefix: "CONFGO"}).Load(&result, f); err != nil {
		t.Error(err)
		t.Fail()
	} else if !reflect.DeepEqual(result, config) {
		t.Errorf("result should equal to original configuration")
	}
}

func TestOverwritetestConfigurationWithEnvironmentWithEmptyPrefix(t *testing.T) {
	config := generateDefaultConfig()
	f := generateFile(t, config)
	defer os.Remove(f)

	os.Setenv("APPNAME", "config2")
	os.Setenv("HOSTS", "- http://example.org\n- http://hello.world")
	os.Setenv("DB_NAME", "db_name")
	defer os.Setenv("APPNAME", "")
	defer os.Setenv("HOSTS", "")
	defer os.Setenv("DB_NAME", "")

	config.APPName = "config2"
	config.Hosts = []string{"http://example.org", "http://hello.world"}
	config.DB.Name = "db_name"

	var result testConfig
	if err := Load(&result, f); err != nil {
		t.Error(err)
		t.Fail()
	} else if !reflect.DeepEqual(result, config) {
		t.Errorf("result should equal to original configuration")
	}
}

func TestOverwritetestConfigurationWithEnvironmentWithCustomKey(t *testing.T) {
	config := generateDefaultConfig()
	f := generateFile(t, config)
	defer os.Remove(f)

	os.Setenv("DBPassword", "mypass")
	defer os.Setenv("DBPassword", "")

	config.DB.Password = "mypass"

	var result testConfig
	if err := Load(&result, f); err != nil {
		t.Error(err)
		t.Fail()
	} else if !reflect.DeepEqual(result, config) {
		t.Log(result.DB.Password)
		t.Errorf("result should equal to original configuration")
	}
}

func TestOverwritetestConfigurationWithArgument(t *testing.T) {
	config := generateDefaultConfig()
	f := generateFile(t, config)
	defer os.Remove(f)

	os.Args = append(os.Args, "-n", "3", "-b", "-h", "this is a help infomation")
	config.Number = 3
	config.Boolean = true
	config.Help = "this is a help infomation"

	var result testConfig
	if err := Load(&result, f); err != nil {
		t.Error(err)
		t.Fail()
	} else if !reflect.DeepEqual(result, config) {
		t.Log(result.Number)
		t.Errorf("result should equal to original configuration")
	}
}
