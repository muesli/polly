import Ember from 'ember';

export default Ember.Service.extend({
  session: Ember.inject.service('session'),
  store: Ember.inject.service('store'),

  load() {
    let userId = this.get('session.data.authenticated.user_id');
    if (!Ember.isEmpty(userId)) {
      return this.get('store').findRecord('user', userId).then((user) => {
        this.set('user', user);
      });
    } else {
      return Ember.RSVP.resolve();
    }
  }
});
