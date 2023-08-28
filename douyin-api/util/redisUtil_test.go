package util

import "testing"

func TestGetVideoFavoriteAndCommentCount(t *testing.T) {
	type args struct {
		videoID int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{videoID: -10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetVideoFavoriteAndCommentCount(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideoFavoriteAndCommentCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVideoFavoriteAndCommentCount() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetVideoFavoriteAndCommentCount() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getUserAllUserData(t *testing.T) {
	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   int64
		want2   int64
		want3   int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{userID: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, err := getUserAllUserData(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserAllUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getUserAllUserData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getUserAllUserData() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("getUserAllUserData() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("getUserAllUserData() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
