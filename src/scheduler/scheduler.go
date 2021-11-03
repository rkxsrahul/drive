package scheduler

import (
	"log"
	"strconv"
	"time"

	"git.xenonstack.com/util/drive-portal/config"

	"github.com/robfig/cron"
	"go.uber.org/zap"
)

var Cron *cron.Cron

// Start is a function to start cronjobs
func Start() {

	log.Println("Scheduler Started...")
	AnswerToResultToPool()
	DrivePoolRank()
	DrivePoolWeek()
	Cron := cron.New()
	Cron.AddFunc("0 0 * * *", AnswerToResultToPool)
	Cron.AddFunc("0 1 * * 0", DrivePoolRank)
	Cron.AddFunc("0 2 */2 * 0", DrivePoolWeek)
	Cron.Start()
}

func AnswerToResultToPool() {
	db := config.DB
	// create result table if not exists
	err := db.Exec("CREATE TABLE if not exists results (like answers INCLUDING ALL);").Error
	if err != nil {
		zap.S().Error("scheduler 1....1....", err)
		return
	}
	// create foreign key if not exists
	db.Exec("alter table results add constraint  result_foreign_key foreign key (drive_id,email) references drive_users(drive_id,user_email) on delete cascade;")
	now := strconv.FormatInt(time.Now().Unix(), 10)
	// insert data in result table
	err = db.Exec("insert into results  (select * from answers where drive_id not in (select drive_id from user_sessions where expire>" + now + ")) ON CONFLICT ON CONSTRAINT results_pkey DO update set marked_id=excluded.marked_id, answer_id=excluded.answer_id, time=excluded.time, ques_id=excluded.ques_id, created_at=excluded.created_at, updated_at=excluded.updated_at ;").Error
	if err != nil {
		zap.S().Error("scheduler 1....2...", err)
		return
	}
	// delete data from answers table
	err = db.Exec("delete from answers where drive_id not in (select drive_id from user_sessions where expire>" + now + ")").Error
	if err != nil {
		zap.S().Error("scheduler 1....3...", err)
		return
	}
	// create schema
	err = db.Exec("CREATE SCHEMA IF NOT EXISTS results").Error
	if err != nil {
		zap.S().Error("scheduler 1....4...", err)
		return
	}
	// create table pool
	err = db.Exec("create table if not exists results.pool_analytical(email varchar NOT NULL,drive_id int8 NOT NULL,pool_id varchar NOT NULL,pool_total_question int NOT NULL,pool_total_attempted int NOT NULL,pool_total_correct int NOT NULL,pool_percentage numeric NOT NULL, pool_accuracy numeric NOT null, time_taken int DEFAULT 0, CONSTRAINT pool_analytical_pkeys PRIMARY KEY (drive_id, email, pool_id));").Error
	if err != nil {
		zap.S().Error("scheduler 1....5...", err)
		return
	}

	err = db.Exec("ALTER TABLE results.pool_analytical ADD CONSTRAINT pool_analytical_pkeys PRIMARY KEY (drive_id, email, pool_id);").Error
	zap.S().Error("scheduler 1....5.1...", err)
	err = db.Exec("create table if not exists results.drive_analytical(email varchar NOT NULL,drive_id int8 NOT NULL,accuracy numeric NOT NULL, CONSTRAINT drive_analytical_pkey PRIMARY KEY (drive_id, email));").Error
	if err != nil {
		zap.S().Error("scheduler 1....6...", err)
		return
	}
	// insert data
	err = db.Exec("insert into results.pool_analytical select email,drive_id,pool_id,count(ques_id) as pool_total,count(marked_id) filter (where marked_id<>0) as pool_total_attempted,count(marked_id) filter (where marked_id=answer_id) as pool_total_correct, ((count(marked_id) filter (where marked_id=answer_id))::decimal/(greatest(count(marked_id),1))::decimal * 100) as pool_percentage, ((count(marked_id) filter (where marked_id=answer_id))::decimal/(greatest(count(marked_id) filter (where marked_id<>0),1))::decimal * 100)  as pool_accuracy from results group by email,drive_id,pool_id order by email ON CONFLICT ON CONSTRAINT pool_analytical_pkeys DO update set pool_total_question=excluded.pool_total_question, pool_total_attempted=excluded.pool_total_attempted, pool_total_correct=excluded.pool_total_correct, pool_percentage=excluded.pool_percentage, pool_accuracy=excluded.pool_accuracy;").Error
	if err != nil {
		zap.S().Error("scheduler 1....7...", err)
		return
	}
	db.Exec("delete from results.pool_analytical where pool_id not in (select id from pools);")
	db.Exec("delete from results.pool_analytical where drive_id not in (select id from drives);")
	err = db.Exec("insert into results.drive_analytical select email,drive_id,(sum(pool_total_correct)::decimal/greatest(sum(pool_total_attempted),1)*100) as accuracy from results.pool_analytical group by email,drive_id order by accuracy desc ON CONFLICT ON CONSTRAINT drive_analytical_pkey DO update set accuracy=excluded.accuracy;").Error
	if err != nil {
		zap.S().Error("scheduler 1....8...", err)
		return
	}
	db.Exec("delete from results.drive_analytical where drive_id not in (select id from drives);")
}

func DrivePoolRank() {
	db := config.DB
	//create tables if not exist
	err := db.Exec("create table if not exists results.drive_rank(drive_id int8 NOT NULL, ranking int8 NOT NULL, calculated_at timestamptz NOT NULL DEFAULT now());").Error
	if err != nil {
		zap.S().Error("scheduler 2....1....", err)
		return
	}
	err = db.Exec("create table if not exists results.pool_rank(pool_id varchar NOT NULL,ranking int8 NOT NULL, calculated_at timestamptz NOT NULL DEFAULT now());").Error
	if err != nil {
		zap.S().Error("scheduler 2....2....", err)
		return
	}
	//create index and cluster
	err = db.Exec("create index if not exists drive_rank_idx on results.drive_rank (drive_id asc);").Error
	if err != nil {
		zap.S().Error("scheduler 2....3....", err)
		return
	}
	err = db.Exec("create index if not exists pool_rank_idx on results.pool_rank (pool_id asc);").Error
	if err != nil {
		zap.S().Error("scheduler 2....4....", err)
		return
	}
	// insert data
	err = db.Exec("insert into results.drive_rank select drive_id, rank() over (order by (sum(pool_total_correct)::decimal/(greatest(sum(pool_total_attempted),1))::decimal*100) DESC) from results.pool_analytical group by drive_id;").Error
	if err != nil {
		zap.S().Error("scheduler 2....5....", err)
		return
	}
	err = db.Exec("insert into results.pool_rank select pool_id, rank() over (order by (sum(pool_total_correct)::decimal/(greatest(sum(pool_total_attempted),1))::decimal*100) DESC) from results.pool_analytical group by pool_id;").Error
	if err != nil {
		zap.S().Error("scheduler 2....6....", err)
		return
	}
}

func DrivePoolWeek() {
	db := config.DB
	//create tables if not exist
	err := db.Exec("create table if not exists results.drive_week(drive_id int8 NOT NULL, total int8 NOT NULL, twenty int8 NOT NULL, forty int8 NOT NULL, sixty int8 NOT NULL, eighty int8 NOT NULL, hundred int8 NOT NULL, calculated_at timestamptz NOT NULL DEFAULT now());").Error
	if err != nil {
		zap.S().Error("scheduler 3....6....", err)
		return
	}
	err = db.Exec("create table if not exists results.pool_week(pool_id varchar NOT NULL, total int8 NOT NULL, twenty int8 NOT NULL, forty int8 NOT NULL, sixty int8 NOT NULL, eighty int8 NOT NULL, hundred int8 NOT NULL, calculated_at timestamptz NOT NULL DEFAULT now());").Error
	if err != nil {
		zap.S().Error("scheduler 3....1....", err)
		return
	}
	//create indexes
	err = db.Exec("create index if not exists drive_week_idx on results.drive_week (drive_id asc);").Error
	if err != nil {
		zap.S().Error("scheduler 3....2....", err)
		return
	}
	err = db.Exec("create index if not exists pool_week_idx on results.pool_week (pool_id asc);").Error
	if err != nil {
		zap.S().Error("scheduler 3....3....", err)
		return
	}
	//insert data
	err = db.Exec("insert into results.drive_week select drive_id,count(distinct email) as total, count(distinct email) filter (where accuracy>=0 and accuracy <20), count(distinct email) filter (where accuracy>=20 and accuracy<40), count(distinct email) filter (where accuracy>=40 and accuracy<60), count(distinct email) filter (where accuracy>=60 and accuracy<80), count(distinct email) filter (where accuracy>=80 and accuracy<=100) from results.drive_analytical group by drive_id;").Error
	if err != nil {
		zap.S().Error("scheduler 3....4....", err)
		return
	}
	err = db.Exec("insert into results.pool_week select pool_id,count(distinct email) as total, count(distinct email) filter (where pool_accuracy>=0 and pool_accuracy <20), count(distinct email) filter (where pool_accuracy>=20 and pool_accuracy<40), count(distinct email) filter (where pool_accuracy>=40 and pool_accuracy<60), count(distinct email) filter (where pool_accuracy>=60 and pool_accuracy<80), count(distinct email) filter (where pool_accuracy>=80 and pool_accuracy<=100) from results.pool_analytical group by pool_id;").Error
	if err != nil {
		zap.S().Error("scheduler 3....5....", err)
		return
	}
}
