import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    activate: function() {
        this._super();
        Ember.$('grantsRunning').button('toggle');
    },

    model() {
       return Ember.RSVP.hash({
         smallGrantProposals: this.store.query('proposal', {
             granttype: 'small'
         }),
         largeGrantProposals: this.store.query('proposal', {
             granttype: 'large'
         })
       });
     },

     setupController(controller, models) {
       controller.set('smallGrantProposals', models.smallGrantProposals);
       controller.set('largeGrantProposals', models.largeGrantProposals);
     }
});
