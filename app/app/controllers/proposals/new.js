import Ember from 'ember';
import moment from 'moment';

export default Ember.Controller.extend({
  title: "",
  description: "",
  activities: "",
  contact: "",
  recipient: "",
  recipient2: "",
  value: "",
  responseMessage: "",
  errorMessage: "",
  startdate: new Date(),
  maxmicrobudget: 0,
  maxvalue: 0,

  isMicroBudget: Ember.computed('value', function() {
      const max = this.get('maxmicrobudget');
      return max === 0 || this.value === '' || parseInt(this.value) <= max;
  }),

  maxBudget: Ember.computed('startdate', 'maxmicrobudget', 'maxvalue', function() {
      this.store.query('budget', {
          month: moment(this.get('startdate')).add(14, 'd').toDate().getMonth() + 1
      }).then((budget) => {
          this.set('maxmicrobudget', budget.objectAt(0).get('value'));
          this.set('maxvalue', budget.objectAt(0).get('maxvalue'));
      });

      return this.get('maxmicrobudget');
  }),

  minimumProposalStartDate: Ember.computed(function() {
      var date = moment().add(1, 'd').toDate();
      return moment(date).format('YYYY/MM/DD');
  }),

  isValid: Ember.computed('recipient', 'recipient2', 'contact', 'title', 'description', 'value', 'startdate', function() {
      return this.title.length > 0 && this.description.length > 0 &&
             this.activities.length > 0 && this.contact.length > 0 &&
             this.recipient.length > 0 && (this.value <= this.maxmicrobudget || this.recipient2.length > 0) &&
             parseInt(this.value) > 0 && parseInt(this.value) <= this.maxvalue &&
             this.startdate.getFullYear() > 0;
  }),
  isDisabled: Ember.computed.not('isValid'),

  actions: {
    createProposal() {
      this.set('errorMessage', '');
      this.set('responseMessage', '');
      this.set('progressMessage', `Creating proposal...`);

      const title = this.get('title');
      const description = this.get('description');
      const activities = this.get('activities');
      const contact = this.get('contact');
      const recipient = this.get('recipient');
      const recipient2 = this.get('recipient2');
      const value = this.get('value');
      const startdate = this.get('startdate');
      const newProposal = this.store.createRecord('proposal', { title: title, description: description, activities: activities, contact: contact, recipient: recipient, recipient2: recipient2, value: value, starts: startdate });
      newProposal.save().then(
        (/*proposal*/) => {
          this.set('responseMessage', `Your proposal is now awaiting moderation. Thank you!`);
          this.set('progressMessage', '');
          this.set('title', '');
          this.set('description', '');
          this.set('activities', '');
          this.set('contact', '');
          this.set('recipient', '');
          this.set('recipient2', '');
          this.set('value', '');
          this.set('startdate', new Date());
        },
        error => {
          this.set('errorMessage', `Failed adding your proposal: ` + error);
          this.set('progressMessage', '');
        }
      );
    }
  }
});
