import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),
  errorMessage: "",
  signupemail: "",
  signuptoken: "",
  password: "",
  confirmPassword: "",

  isValid: Ember.computed('password', 'confirmPassword', function() {
    return this.get('password').length > 6 && this.get('password') === this.get('confirmPassword');
  }),
  isDisabled: Ember.computed.not('isValid'),

  actions: {
    signup() {
      this.set('errorMessage', '');
      this.set('responseMessage', '');
      this.set('progressMessage', `Creating account for ${this.get('signupemail')}...`);

      const token = this.get('signuptoken');
      const password = this.get('password');
      var options = {token: token, password: password};

      this.get('session').authenticate('authenticator:custom', options).catch((reason) => {
          this.set('errorMessage', `Failed creating account for ${this.get('signupemail')}: ` + reason);
          this.set('progressMessage', '');
      });
    }
  }
});
