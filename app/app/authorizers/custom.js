import Ember from 'ember';
import Base from 'ember-simple-auth/authorizers/base';

export default Base.extend({
    authorize(data, block) {
        const accessToken = data['token'];
        if (!Ember.isEmpty(accessToken)) {
            block('Authorization', `Bearer ${accessToken}`);
        }
    }
});
