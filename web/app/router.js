import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('folder', {resetNamespace: true, path: '/folders/:folder_id'}, function () {
    this.route('page', {resetNamespace: true, path: '/pages/:page_id'}, function () {
      this.route('show', {path: ''});
      this.route('edit');
    });
  });
  this.route('not-found', {path: '/*path'});
  this.route('pages', {path: '/pages'});
});

export default Router;
