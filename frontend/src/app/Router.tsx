import { HomePage } from '../pages/home/index';
import { ModeratorPage } from '../pages/moderator/index';
import { ModeratorLoginPage } from '../pages/moderator/index';
import { Router, Route, Navigate } from '@solidjs/router';
import { StatusPage } from '../pages/status/index';
import { render } from 'solid-js/web';

const root = document.getElementById('root');
render(
  () => (
    <Router>
      <Route path='/moderator'>
        <Route
          path='/'
          component={ModeratorPage}
        />
        <Route
          path='/login'
          component={ModeratorLoginPage}
        />
      </Route>
      <Route
        path='/status'
        component={StatusPage}
      />
      <Route
        path='/*'
        component={() => (
          <Navigate
            href='/'
            state
          />
        )}
      />
      <Route
        path='/'
        component={HomePage}
      />
    </Router>
  ),
  root!,
);
