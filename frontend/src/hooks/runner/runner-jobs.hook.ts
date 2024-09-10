import { useEffect, useState } from 'react';
import { fetchRunnerJobs } from '../../services/runners.service';
import { IRunnerJobs } from '../../interfaces/runner.interface';
import { IStructuredError } from '../../interfaces/structured-error.interface';

interface IUseRunnerJobs {
  runnerJobs?: IRunnerJobs;
  runnerJobsError?: IStructuredError;
  runnerJobsLoading: boolean;
}

export function useRunnerJobs(url: string): IUseRunnerJobs {
  const [runnerJobs, setRunnerJobs] = useState<IRunnerJobs | undefined>();
  const [runnerJobsError, setRunnerJobsError] = useState<IStructuredError | undefined>();
  const [runnerJobsLoading, setRunnerJobsLoading] = useState(true);

  useEffect(() => {
    const runFetchRunner = async (): Promise<void> => {
      setRunnerJobsLoading(true);

      await fetchRunnerJobs(url)
        .then((response) => {
          setRunnerJobs(response);
        })
        .catch((error: IStructuredError) => {
          setRunnerJobsError(error);
        })
        .finally(() => {
          setRunnerJobsLoading(false);
        });
    };

    runFetchRunner();
  }, [url]);

  return { runnerJobs, runnerJobsError, runnerJobsLoading };
}
