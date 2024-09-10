import { IJob } from './job.interface';
import { ICommit } from './commit.interface';
import { IBuild } from './build.interface';
import { IRepo } from './repo.interface';

export interface IRunner {
  architecture: string;
  created_at: string;
  deleted_at?: string;
  etag: string;
  id: string;
  labels: string[];
  legal_entity_id: string;
  name: string;
  enabled: boolean;
  operating_system: string;
  software_version: string;
  supported_job_types: string[];
  updated_at: string;
  url: string;
  runners_jobs_url: string;
}

export interface IRunnerJobResults {
  job: IJob;
  build: IBuild;
  commit: ICommit;
  repo: IRepo;
}

export interface IRunnerJobs {
  id: string;
  url: string;
  created_at: string;
  runner_job_results: IRunnerJobResults[];
}
