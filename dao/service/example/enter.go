package example

type ExampleService struct {
	SysUserService
	TaskConfigService
	SubTaskConfigService
	MysqlClusterService
	MysqlDBService
	MysqlInstanceService
	MysqlSessionService
	PgSessionService
}
