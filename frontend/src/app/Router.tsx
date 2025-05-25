import { render } from 'solid-js/web';
import { Router, Route, Navigate } from '@solidjs/router';
import { HomePage } from '../pages/home/index';
import { ModeratorPage } from '../pages/moderator/index';
import { StatusPage } from '../pages/status/index';

const root = document.getElementById('root');
render(
  () => (
    <Router>
      <Route
        path='/home'
        component={HomePage}
      />
      <Route
        path='/moderator'
        component={ModeratorPage}
      />
      <Route
        path='/status/:type/:id'
        component={StatusPage}
      />
      <Route
        path='/*'
        component={() => (
          <Navigate
            href='/home'
            state
          />
        )}
      />
    </Router>
  ),
  root!,
);
