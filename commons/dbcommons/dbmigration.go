package dbcommons

import (
	"bufio"
	"embed"
	"fmt"
	"im-server/commons/tools"
	"sort"
	"strings"
)

//go:embed sqls/*
var sqlFs embed.FS

const (
	JimDbVersionKey = "jimdb_versaion"
)

func Upgrade() {
	var currVersion int64 = 0
	dao := GlobalConfDao{}
	conf, err := dao.FindByKey(JimDbVersionKey)
	if err == nil {
		ver, err := tools.String2Int64(conf.ConfValue)
		if err == nil && ver > 0 {
			currVersion = ver
		}
	} else {
		dao.Create(GlobalConfDao{
			ConfKey:   JimDbVersionKey,
			ConfValue: "20240716",
		})
	}
	fmt.Println("[DbMigration]current version:", currVersion)
	sqlFiles, err := sqlFs.ReadDir("sqls")
	if err == nil {
		neededVers := []int64{}
		for _, sqlFile := range sqlFiles {
			fileName := sqlFile.Name()
			if len(fileName) == 12 {
				fileName = fileName[:8]
			}
			ver, err := tools.String2Int64(fileName)
			if err == nil && ver > 0 {
				neededVers = append(neededVers, ver)
			}
		}
		//sort
		sort.Slice(neededVers, func(i, j int) bool {
			return neededVers[i] < neededVers[j]
		})
		for _, ver := range neededVers {
			if ver > currVersion {
				sqlFileName := fmt.Sprintf("sqls/%d.sql", ver)
				fmt.Println("[DbMigration]start to execute sql file:", sqlFileName)
				err := executeSqlFile(sqlFileName)
				if err == nil {
					fmt.Println("[DbMigration]execute sql file success:", sqlFileName)
					dao.UpdateValue(JimDbVersionKey, fmt.Sprintf("%d", ver))
				}
			}
		}
	}
}

func executeSqlFile(fileName string) error {
	sqlFile, err := sqlFs.Open(fileName)
	if err != nil {
		fmt.Println("[DbMigration_Err]Read sql file err:", err, "file_name:", fileName)
		return err
	}
	defer sqlFile.Close()

	scanner := bufio.NewScanner(sqlFile)
	var queryBuilder strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "--") {
			continue
		}
		queryBuilder.WriteString(line)
		queryBuilder.WriteByte(' ')
		if strings.HasSuffix(line, ";") {
			query := strings.TrimSpace(queryBuilder.String())
			if query != "" {
				if err := GetDb().Exec(query).Error; err != nil {
					fmt.Println("[DbMigration_Err]Execute sql error:", err, query)
				}
			}
			queryBuilder.Reset()
		}
	}
	return nil
}
