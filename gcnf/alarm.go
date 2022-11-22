package gcnf

import "github.com/layasugar/laya/core/constants"

func AlarmType() string {
	if IsSet(constants.KEY_APPALARMTYPE) {
		return GetString(constants.KEY_APPALARMTYPE)
	}
	return defaultNullString
}

func AlarmKey() string {
	if IsSet(constants.KEY_APPALARMKEY) {
		return GetString(constants.KEY_APPALARMKEY)
	}
	return defaultNullString
}

func AlarmHost() string {
	if IsSet(constants.KEY_APPALARMADDR) {
		return GetString(constants.KEY_APPALARMADDR)
	}
	return defaultNullString
}
