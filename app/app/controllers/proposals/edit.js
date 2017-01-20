import Ember from 'ember';
import moment from 'moment';

export default Ember.Controller.extend({
  responseMessage: "",
  errorMessage: "",
  maxmicrobudget: 0,

  title: Ember.computed(function() {
      return this.get('proposal').get('title');
  }),
  description: Ember.computed(function() {
      return this.get('proposal').get('description');
  }),
  recipient: Ember.computed(function() {
      return this.get('proposal').get('recipient');
  }),
  value: Ember.computed(function() {
      return this.get('proposal').get('value');
  }),
  startdate: Ember.computed(function() {
      return this.get('proposal').get('starts');
  }),

  maxBudget: Ember.computed('startdate', 'maxmicrobudget', function() {
      this.store.query('budget', {
          month: moment(this.get('startdate')).add(14, 'd').toDate().getMonth() + 1
      }).then((budget) => {
          this.set('maxmicrobudget', budget.objectAt(0).get('value'));
      });

      return this.get('maxmicrobudget');
  }),

  minimumProposalStartDate: Ember.computed(function() {
      var date = moment().add(1, 'd').toDate();
      return moment(date).format('YYYY/MM/DD');
  }),

  isValid: Ember.computed('recipient', 'title', 'description', 'value', 'startdate', function() {
      const title = this.get('title');
      const description = this.get('description');
      const recipient = this.get('recipient');
      const value = this.get('value');
      const startdate = this.get('startdate');

      return title.length > 0 && description.length > 0 &&
             recipient.length > 0 && parseInt(value) > 0 &&
             startdate.getFullYear() > 0;
  }),
  isDisabled: Ember.computed.not('isValid'),

  actions: {
    saveProposal() {
      this.set('errorMessage', '');
      this.set('responseMessage', '');
      this.set('progressMessage', `Saving proposal...`);

      const title = this.get('title');
      const description = this.get('description');
      const recipient = this.get('recipient');
      const value = this.get('value');
      const startdate = this.get('startdate');

      var proposal = this.get('proposal');

      proposal.set('title', title);
      proposal.set('description', description);
      proposal.set('recipient', recipient);
      proposal.set('value', value);
      proposal.set('starts', startdate);

      proposal.save().then(
        (/*proposal*/) => {
          this.set('responseMessage', `Your proposal has been updated. Thank you!`);
          this.set('progressMessage', '');
/*          this.set('title', '');
          this.set('description', '');
          this.set('recipient', '');
          this.set('value', '');
          this.set('startdate', new Date());*/
        },
        error => {
          this.set('errorMessage', `Failed updating your proposal: ` + error);
          this.set('progressMessage', '');
        }
      );
    }
  }
});
