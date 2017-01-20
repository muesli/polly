import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('about');
  this.route('admin');
  this.route('proposals', function() {
    this.route('new');
    this.route('show', { path: '/:proposal_id' });
    this.route('all');
    this.route('finished');
    this.route('accepted');
    this.route('rejected');
    this.route('edit', { path: '/edit/:proposal_id' });
    this.route('own');
  });
  this.route('login');
  this.route('signup', { path: '/signup/:token' });
});

export default Router;
