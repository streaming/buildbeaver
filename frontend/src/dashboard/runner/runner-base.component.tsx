import React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { RunnerEdit } from './runner-edit.component';
import { RunnerView } from './runner-view.component';
import { RunnerJobs } from './runner-jobs.component';

export function RunnerBase(): JSX.Element {
  return (
    <Routes>
      <Route path="/" element={<RunnerView />} />
      <Route path="/edit" element={<RunnerEdit />} />
      <Route path="/jobs" element={<RunnerJobs />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
