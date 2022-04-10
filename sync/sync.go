package sync

func Sync(refs []string) {
	for _, ref := range refs {
		syncRef(ref)
	}
}

func syncRef(ref string) {

}

func isConflictedMergeRef(ref string) bool {
	return false //TODO: Если случился конфликт при мерже, то создается ветка со специальным именем и она потом обрабатывается специальным образом
}
