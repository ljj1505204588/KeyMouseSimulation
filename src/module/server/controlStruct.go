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
	case enum.Free:
		if change == enum.RecordPause {
			return fmt.Errorf(language.ErrorFreeToRecordPause)
		} else if change == enum.PlaybackPause {
			return fmt.Errorf(language.ErrorFreeToPlaybackPause)
		}
	// -------------- 回放 --------------
	case enum.Playback:
		if change == enum.Recording || change == enum.RecordPause {
			return fmt.Errorf(language.ErrorPlaybackToRecordOrRecordPause)
		}
	// -------------- 回放暂停 --------------
	case enum.PlaybackPause:
		if change == enum.Recording || change == enum.RecordPause {
			return fmt.Errorf(language.ErrorPlaybackPauseToRecordOrRecordPause)
		}
	// -------------- 记录 --------------
	case enum.Recording:
		if change == enum.Playback || change == enum.PlaybackPause {
			return fmt.Errorf(language.ErrorRecordToPlaybackOrPlaybackPause)
		}
	// -------------- 记录暂停 --------------
	case enum.RecordPause:
		if change == enum.Playback || change == enum.PlaybackPause {
			return fmt.Errorf(language.ErrorRecordPauseToPlaybackOrPlaybackPause)
		}
	}

	t.statusEnum = change
	return nil
}

// 获取暂停后状态
func (t *statusT) getAfterPauseStatus() (status enum.Status, err error) {
	switch t.statusEnum {
	case enum.Playback:
		return enum.PlaybackPause, nil
	case enum.PlaybackPause:
		return enum.PlaybackPause, nil
	case enum.Recording:
		return enum.RecordPause, nil
	case enum.RecordPause:
		return enum.RecordPause, nil
	default:
		return t.statusEnum, errors.New(language.ErrorPauseFail)
	}
}
