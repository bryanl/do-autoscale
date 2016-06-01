import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('templates');
  this.route('groups', function() {
    this.route('new');
    this.route('show', { path: '/:group_id'});
    this.route('policy', { path: '/:group_id/policy'});
    this.route('delete', { path: '/:group_id/delete'});
  });
});

export default Router;
