import { JSX } from 'solid-js';

import { Router, Route, Navigate } from '@solidjs/router';

import HomePage from './pages/HomePage';
import AddThingPage from './pages/AddThingPage';
import ThingStatusPage from './pages/ThingStatusPage';
import {
  HOME_ROUTE,
  ADD_THING_ROUTE,
  THING_STATUS_ROUTE,
} from './utils/consts';

const AppRouter = (): JSX.Element => {
  return (
    <Router>
      <Route
        path={ADD_THING_ROUTE}
        component={() => <AddThingPage />}
      />
      <Route
        path={THING_STATUS_ROUTE}
        component={() => <ThingStatusPage />}
      />
      <Route
        path={HOME_ROUTE}
        component={() => <HomePage />}
      />
      <Route
        path='/*'
        component={() => (
          <Navigate
            href={HOME_ROUTE}
            state
          />
        )}
      />
    </Router>
  );
};

export default AppRouter;
