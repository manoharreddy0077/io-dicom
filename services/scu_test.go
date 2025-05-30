package services

import (
	"testing"

	"github.com/manoharreddy0077/io-dicom/dictionary/tags"
	"github.com/manoharreddy0077/io-dicom/media"
	"github.com/manoharreddy0077/io-dicom/network"
	"github.com/manoharreddy0077/io-dicom/network/dicomstatus"
	"github.com/manoharreddy0077/io-dicom/utils"
)

func Test_scu_EchoSCU(t *testing.T) {
	_, testSCP := StartSCP(t, 1040)

	testSCP.OnAssociationRequest(func(request network.AAssociationRQ) bool {
		return request.GetCalledAE() == "TEST_SCP"
	})

	media.InitDict()

	type fields struct {
		destination *network.Destination
	}
	type args struct {
		timeout int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should have C-Echo Success",
			fields: fields{
				destination: &network.Destination{
					Name:      "Test Destination",
					CalledAE:  "TEST_SCP",
					CallingAE: "TEST_SCU",
					HostName:  "localhost",
					Port:      1040,
					IsCFind:   false,
					IsCMove:   false,
					IsCStore:  false,
					IsTLS:     false,
				},
			},
			args: args{
				timeout: 0,
			},
			wantErr: false,
		},
		{
			name: "Should not have C-Echo Success",
			fields: fields{
				destination: &network.Destination{
					Name:      "Test Destination",
					CalledAE:  "TEST_SCP2",
					CallingAE: "TEST_SCU",
					HostName:  "localhost",
					Port:      1040,
					IsCFind:   false,
					IsCMove:   false,
					IsCStore:  false,
					IsTLS:     false,
				},
			},
			args: args{
				timeout: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewSCU(tt.fields.destination)
			if err := d.EchoSCU(tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("scu.EchoSCU() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_scu_FindSCU(t *testing.T) {
	_, testSCP := StartSCP(t, 1041)

	testSCP.OnAssociationRequest(func(request network.AAssociationRQ) bool {
		return request.GetCalledAE() == "TEST_SCP"
	})

	testSCP.OnCFindRequest(func(request network.AAssociationRQ, findLevel string, data media.DcmObj) ([]media.DcmObj, uint16) {
		return make([]media.DcmObj, 0), dicomstatus.Success
	})

	media.InitDict()

	type fields struct {
		destination *network.Destination
	}
	type args struct {
		Query   media.DcmObj
		timeout int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint16
		wantErr bool
	}{
		{
			name: "Should C-Find All",
			fields: fields{
				destination: &network.Destination{
					Name:      "Test Destination",
					CalledAE:  "TEST_SCP",
					CallingAE: "TEST_SCU",
					HostName:  "localhost",
					Port:      1041,
					IsCFind:   true,
					IsCMove:   true,
					IsCStore:  true,
					IsTLS:     false,
				},
			},
			args: args{
				Query:   utils.DefaultCFindRequest(),
				timeout: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.Query.WriteString(tags.StudyDate, "20150617")
			d := NewSCU(tt.fields.destination)
			d.SetOnCFindResult(func(result media.DcmObj) {
				result.DumpTags()
			})

			_, status, err := d.FindSCU(tt.args.Query, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("scu.FindSCU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if status != tt.want {
				t.Errorf("scu.FindSCU() = %v, want %v", status, tt.want)
			}
		})
	}
}

func Test_scu_StoreSCU(t *testing.T) {
	_, testSCP := StartSCP(t, 1042)

	testSCP.OnAssociationRequest(func(request network.AAssociationRQ) bool {
		return request.GetCalledAE() == "TEST_SCP"
	})

	testSCP.OnCStoreRequest(func(request network.AAssociationRQ, data media.DcmObj) uint16 {
		data.DumpTags()
		return dicomstatus.Success
	})

	media.InitDict()

	type fields struct {
		destination *network.Destination
	}
	type args struct {
		FileName string
		timeout  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should store DICOM file",
			fields: fields{
				destination: &network.Destination{
					Name:      "Test Destination",
					CalledAE:  "TEST_SCP",
					CallingAE: "TEST_SCU",
					HostName:  "localhost",
					Port:      1042,
					IsCFind:   true,
					IsCMove:   true,
					IsCStore:  true,
					IsTLS:     false,
				},
			},
			args: args{
				FileName: "../samples/test.dcm",
				timeout:  0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewSCU(tt.fields.destination)
			if err := d.StoreSCU(tt.args.FileName, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("scu.StoreSCU() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func StartSCP(t testing.TB, port int) (func(t testing.TB), SCP) {
	testSCP := NewSCP(port)
	go func() {
		if err := testSCP.Start(); err != nil {
			panic(err)
		}
	}()

	return func(t testing.TB) {
		if err := testSCP.Stop(); err != nil {
			panic(err)
		}
	}, testSCP
}
