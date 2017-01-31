import Ember from 'ember';
import moment from 'moment';

export default Ember.Controller.extend({
    maxmicrobudget: 0,
    periodend: "",
    largegrantperiodend: "",

    periodEnd: Ember.computed('periodend', function() {
        return this.periodend;
    }),

    largeGrantPeriodEnd: Ember.computed('largegrantperiodend', function() {
        return this.largegrantperiodend;
    }),

    maxBudget: Ember.computed('maxmicrobudget', function() {
        this.store.query('budget', {
            month: moment().add(14, 'd').toDate().getMonth() + 1
        }).then((budget) => {
            this.set('maxmicrobudget', budget.objectAt(0).get('value'));
            this.set('periodend', budget.objectAt(0).get('period_end'));
            this.set('largegrantperiodend', budget.objectAt(0).get('large_grant_period_end'));
        });

        return this.get('maxmicrobudget');
    })
});
