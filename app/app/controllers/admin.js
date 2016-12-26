import Ember from 'ember';

export default Ember.Controller.extend({
  emailAddress: "",
  name: "",
  responseMessage: "",
  errorMessage: "",

  activeUserCount: Ember.computed('users.@each.activated', function() {
      var count = 0;
      this.get('users').forEach(function(entry) {
          if (entry.get('activated')) {
              count++;
          }
      });
      return count;
  }),
  inactiveUserCount: Ember.computed('users.@each.activated', function() {
      var count = 0;
      this.get('users').forEach(function(entry) {
          if (!entry.get('activated')) {
              count++;
          }
      });
      return count;
  }),

  isValid: Ember.computed.match('emailAddress', /^.+@.+\..+$/),
  isDisabled: Ember.computed.not('isValid'),

  actions: {
    saveInvitation() {
      this.set('errorMessage', '');
      this.set('responseMessage', '');
      this.set('progressMessage', `Sending invitation to ${this.get('emailAddress')}...`);

      const email = this.get('emailAddress');
      const newUser = this.store.createRecord('user', { email: email });
      newUser.save().then(
        user => {
          this.set('responseMessage', `An invitation to ${this.get('emailAddress')} has been sent!`);
          this.set('progressMessage', '');
          this.set('emailAddress', '');
          this.set('name', '');
        },
        error => {
          this.set('errorMessage', `Failed sending an invitation to ${this.get('emailAddress')}: ` + error);
          this.set('progressMessage', '');
        }
      );
    }
  }
});
