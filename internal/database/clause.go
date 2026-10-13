package database

import "gorm.io/gorm/clause"

func buildOnConflict(options UpsertOptions) clause.OnConflict {
	onConflict := clause.OnConflict{
		Columns: toClauseColumns(options.ConflictColumns),
	}

	if options.DoNothing {
		onConflict.DoNothing = true
	} else {
		onConflict.DoUpdates = clause.AssignmentColumns(
			options.UpdateColumns,
		)
	}

	return onConflict
}
