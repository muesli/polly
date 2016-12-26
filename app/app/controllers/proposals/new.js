import Ember from 'ember';

export default Ember.Controller.extend({
  title: "",
  description: "",
  recipient: "",
  value: "",
  responseMessage: "",
  errorMessage: "",

  isValid: Ember.computed('recipient', 'title', 'description', 'value', function() {
      return this.title.length > 0 && this.description.length > 0 && this.recipient.length > 0 && parseInt(this.value) > 0;
  }),
  isDisabled: Ember.computed.not('isValid'),

  actions: {
    createProposal() {
      this.set('errorMessage', '');
      this.set('responseMessage', '');
      this.set('progressMessage', `Creating proposal...`);

      const title = this.get('title');
      const description = this.get('description');
      const email = this.get('recipient');
      const value = this.get('value');
      const newUser = this.store.createRecord('proposal', { title: title, description: description, recipient: email, value: value });
      newUser.save().then(
        user => {
          this.set('responseMessage', `Your proposal is now awaiting moderation. Thank you!`);
          this.set('progressMessage', '');
          this.set('title', '');
          this.set('description', '');
          this.set('recipient', '');
          this.set('value', '');
        },
        error => {
          this.set('errorMessage', `Failed adding your proposal: ` + error);
          this.set('progressMessage', '');
        }
      );
    }
  }
});
