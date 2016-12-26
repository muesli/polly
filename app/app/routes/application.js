import Ember from 'ember';
import ApplicationRouteMixin from 'ember-simple-auth/mixins/application-route-mixin';

const { service } = Ember.inject;

export default Ember.Route.extend(ApplicationRouteMixin, {
  session: Ember.inject.service('session'),
  currentUser: service(),

  beforeModel() {
    return this._loadCurrentUser();
  },

  sessionAuthenticated() {
    this._super(...arguments);
    this._loadCurrentUser().catch(() => this.get('session').invalidate());
  },

  sessionInvalidated() {
    if (this.get('session.skipRedirectOnInvalidation')) {
      this.set('session.skipRedirectOnInvalidation', false);
    }
    else {
      this._super(...arguments);
    }
  },

  _loadCurrentUser() {
    return this.get('currentUser').load();
  }
});
