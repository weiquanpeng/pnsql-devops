import { request } from '@/utils/request';

const Api = {
  MysqlVmName: '/mysql/vmname',
  PgSession: '/postgresql/session',
};

export function getMysqlVmName() {
  return request.get({
    url: Api.MysqlVmName,
  });
}

export function getPgSession(data: Record<string, any>) {
  return request.post({
    url: Api.PgSession,
    data,
  });
}
