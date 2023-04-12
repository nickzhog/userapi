package filestorage

import (
	"encoding/json"
)

func (r *repository) saveUsersToFile() error {
	err := r.file.Truncate(0)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	_, err = r.file.Seek(0, 0)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	err = json.NewEncoder(r.file).Encode(r.list)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}
