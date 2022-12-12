package server

import (
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/language"
	"errors"
	"fmt"
)

// --------------------------- status ---------------------------

type statusT struct {
	statusEnum enum.Status
}

// 修改状态
func (t *statusT) changeStatus(change enum.Status) (err error) {
	var current = t.statusEnum
	switch current {
	// -------------- 空闲 --------------
	case enum.FREE:
		if change == enum.RECORD_PAUSE {
			return fmt.Errorf(language.ErrorFreeToRecordPause)
		} else if change == enum.PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorFreeToPlaybackPause)
		}
	// -------------- 回放 --------------
	case enum.PLAYBACK:
		if change == enum.RECORDING || change == enum.RECORD_PAUSE {
			return fmt.Errorf(language.ErrorPlaybackToRecordOrRecordPause)
		}
	// -------------- 回放暂停 --------------
	case enum.PLAYBACK_PAUSE:
		if change == enum.RECORDING || change == enum.RECORD_PAUSE {
			return fmt.Errorf(language.ErrorPlaybackPauseToRecordOrRecordPause)
		}
	// -------------- 记录 --------------
	case enum.RECORDING:
		if change == enum.PLAYBACK || change == enum.PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorRecordToPlaybackOrPlaybackPause)
		}
	// -------------- 记录暂停 --------------
	case enum.RECORD_PAUSE:
		if change == enum.PLAYBACK || change == enum.PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorRecordPauseToPlaybackOrPlaybackPause)
		}
	}

	t.statusEnum = change
	return nil
}

// 获取暂停后状态
func (t *statusT) getAfterPauseStatus() (status enum.Status, err error) {
	switch t.statusEnum {
	case enum.PLAYBACK:
		return enum.PLAYBACK_PAUSE, nil
	case enum.PLAYBACK_PAUSE:
		return enum.PLAYBACK_PAUSE, nil
	case enum.RECORDING:
		return enum.RECORD_PAUSE, nil
	case enum.RECORD_PAUSE:
		return enum.RECORD_PAUSE, nil
	default:
		return t.statusEnum, errors.New(language.ErrorPauseFail)
	}
}
