package internal

//New temporary file
// file, err := ioutil.TempFile(".", "sops")
// if err != nil {
// 	logrus.WithError(err).Error("Temporary file for configuration failed")
// }
// defer func() {
// 	logrus.WithError(err).Debug("Defer temporary file")
// 	os.Remove(file.Name())
// }()
// fmt.Println(file.Name()) // For example "dir/prefix054003078"

// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
// if err != nil {
// 	logrus.WithError(err).Error("Init, no filepath")
// 	return
// }
// fmt.Println(dir)

// we use https://age-encryption.org to encrypt/decrypt
// an age key for unit-tests purpose was created with the following command:
// $ age-keygen -o testdata/age.key
// if you need to regenerate it, you'll also need to update its public key here:

/*func TestMain(m *testing.M) {
	log.Println("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")
	os.Exit(exitVal)
}

func TestA(t *testing.T) {
	log.Println("TestA running")

}

func TestB(t *testing.T) {
	log.Println("TestB running")
}

func TestGetFixedValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/fixedvalue" {
			t.Errorf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"value":"fixed"}`))
	}))
	defer server.Close()

	value, _ := GetFixedValue(server.URL)
	if value != "fixed" {
		t.Errorf("Expected 'fixed', got %s", value)
	}
}


func Hash(str, key string) string {
	return str
}

func BenchmarkHash10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		key := "1234567890"
		Hash("6368616e676520746869732070617373", key)
	}
}
func BenchmarkHash20(b *testing.B) {
	for n := 0; n < b.N; n++ {
		key := "12345678901234567890"
		Hash("6368616e676520746869732070617373", key)
	}
}*/

/*func Test_Slow1(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Name     string
		Duration time.Duration
	}{
		{Name: "Case 1", Duration: time.Second * 1},
		{Name: "Case 2", Duration: time.Second * 1},
		{Name: "Case 3", Duration: time.Second * 1},
		{Name: "Case 4", Duration: time.Second * 1},
		{Name: "Case 5", Duration: time.Second * 1},
	}
	for _, c := range cases {
		tc := c
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			t.Logf("%s sleeping..", tc.Name)
			sleepFn(tc.Duration)
			t.Logf("%s slept", tc.Name)
		})
	}
}
func sleepFn(duration time.Duration) {
	time.Sleep(duration)
}
*/
