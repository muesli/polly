import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    model() {
       return Ember.RSVP.hash({
         smallGrantProposals: this.store.query('proposal', {
             granttype: 'small'
         }),
         largeGrantProposals: this.store.query('proposal', {
             granttype: 'large'
         }),
         finishedProposals: this.store.query('proposal', {
             ended: true
         })
       });
     },

     setupController(controller, models) {
       controller.set('smallGrantProposals', models.smallGrantProposals);
       controller.set('largeGrantProposals', models.largeGrantProposals);
       controller.set('finishedProposals', models.finishedProposals);
     }
});
