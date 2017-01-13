import Ember from 'ember';
import moment from 'moment';

export default Ember.Controller.extend({
  title: "",
  description: "",
  recipient: "",
  value: "",
  responseMessage: "",
  errorMessage: "",
  enddate: new Date(),

  minimumProposalEndDate: Ember.computed(function() {
      var date = moment().add(7, 'd').toDate();
      return moment(date).format('YYYY/MM/DD');
  }),

  isValid: Ember.computed('recipient', 'title', 'description', 'value', 'enddate', function() {
      return this.title.length > 0 && this.description.length > 0 &&
             this.recipient.length > 0 && parseInt(this.value) > 0 &&
             this.enddate.getFullYear() > 0;
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
      const enddate = this.get('enddate');
      const newProposal = this.store.createRecord('proposal', { title: title, description: description, recipient: email, value: value, ends: enddate });
      newProposal.save().then(
        (/*proposal*/) => {
          this.set('responseMessage', `Your proposal is now awaiting moderation. Thank you!`);
          this.set('progressMessage', '');
          this.set('title', '');
          this.set('description', '');
          this.set('recipient', '');
          this.set('value', '');
          this.set('enddate', '');
        },
        error => {
          this.set('errorMessage', `Failed adding your proposal: ` + error);
          this.set('progressMessage', '');
        }
      );
    }
  }
});
