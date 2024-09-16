import React from 'react';
import { useRunnerJobsUrl } from '../../hooks/runner-url/runner-url.hook';
import { SimpleContentLoader } from '../../components/content-loaders/simple/simple-content-loader';
import { StructuredError } from '../../components/structured-error/structured-error.component';
import { Button } from '../../components/button/button.component';
import { NavLink } from 'react-router-dom';
import { NotFound } from '../../components/not-found/not-found.component';
import { useRunnerJobs } from '../../hooks/runner/runner-jobs.hook';
import { buildDuration, createdAt, trimRef, trimSha } from '../../utils/build.utils';
import { BuildStatusIndicator } from '../build-status-indicator/build-status-indicator.component';
import { IoCalendarClearOutline, IoStopwatchOutline } from 'react-icons/io5';
import { useAnyLegalEntity } from '../../hooks/any-legal-entity/any-legal-entity.hook';
import { List } from '../../components/list/list.component';
import { Information } from '../../components/information/information.component';
import { getColourForStatus } from '../../utils/build-status.utils';

export function RunnerJobs(): JSX.Element {
  const legalEntity = useAnyLegalEntity();
  const runnerJobsUrl = useRunnerJobsUrl();
  const { runnerJobs, runnerJobsError, runnerJobsLoading } = useRunnerJobs(runnerJobsUrl);

  if (runnerJobsError) {
    return <StructuredError error={runnerJobsError} handleNotFound={true} />;
  }

  if (runnerJobsLoading) {
    return <SimpleContentLoader numberOfRows={1} rowHeight={200} />;
  }

  if (!runnerJobs) {
    return <NotFound />;
  }

  const renderRunnerJobResults = () => {
    if (runnerJobs.runner_job_results == null) {
      return (
        <Information text={'This Runner is yet to run any jobs'} />
      )
    }

    return (
      <List>
        {runnerJobs.runner_job_results.map((result, index: number) => (
          <NavLink
            key={index}
            className={`flex cursor-pointer justify-between text-sm text-gray-600 hover:bg-gray-100`}
            to={`/${legalEntity.type}/${legalEntity.name}/repos/${result.repo.name}/builds/${result.build.name}`}
          >
            <div className="flex min-w-0 grow">
              <div className={`w-[4px] min-w-[4px] bg-${getColourForStatus(result.job.status)}`}></div>
              <div className="flex w-[40%] min-w-[40%] gap-x-2 p-3">
                <div className="flex flex-col items-center gap-y-0.5">
                  <BuildStatusIndicator status={result.job.status} size={20}></BuildStatusIndicator>
                </div>
                <div className="flex min-w-0 flex-col gap-y-1">
                  <div className="... truncate" title={`Build #${result.build.name}`}>
                    {legalEntity.name} / {result.repo.name} <strong>#{result.build.name}</strong>
                  </div>
                  <div className="... truncate" title={`Job #${result.job.name}`}>
                    {result.job.name}
                  </div>
                </div>
              </div>
              <div className="flex min-w-0 grow flex-col gap-y-1 p-3">
                <div className="... truncate" title="Commit message">
                  {result.commit.message}
                </div>
                <div className="text-xs">
                  <span title="Branch">
                    <strong>{trimRef(result.build.ref)}</strong> /{' '}
                  </span>
                  <span title="Sha">{trimSha(result.commit.sha)}</span>
                </div>
              </div>
            </div>
            <div className="flex w-[150px] min-w-[150px] flex-col justify-center gap-y-1 p-3 text-xs">
              <div className="flex items-center gap-x-2" title="Created">
                <IoCalendarClearOutline className="ml-[1px]" size={14} />
                {createdAt(result.job.created_at)}
              </div>
              <div className="flex items-center gap-x-2" title="Duration">
                <IoStopwatchOutline size={16} />
                {buildDuration(result.job.timings)}
              </div>
            </div>
          </NavLink>
        ))}
      </List>
    )
  }

  return (
    <>
      <div className="flex items-center justify-between mb-5">
        <span className="text-lg">Runner Jobs</span>
        <div className="flex gap-x-5">
          <NavLink to="..">
            <Button label="Runner" size="small" />
          </NavLink>
        </div>
      </div>
      {renderRunnerJobResults()}
    </>
  );
}
