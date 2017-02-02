import Ember from 'ember';
import moment from 'moment';

export default Ember.Controller.extend({
    avail_micro_budget: 0,
    periodend: "",
    largegrantperiodend: "",

    periodEnd: Ember.computed('periodend', function() {
        return this.periodend;
    }),

    largeGrantPeriodEnd: Ember.computed('largegrantperiodend', function() {
        return this.largegrantperiodend;
    }),

    availMicroBudget: Ember.computed('avail_micro_budget', function() {
        this.store.query('budget', {
            month: moment().add(14, 'd').toDate().getMonth() + 1
        }).then((budget) => {
            this.set('avail_micro_budget', budget.objectAt(0).get('available_small'));
            this.set('periodend', budget.objectAt(0).get('period_end'));
            this.set('largegrantperiodend', budget.objectAt(0).get('large_grant_period_end'));
        });

        return this.get('avail_micro_budget');
    })
});
