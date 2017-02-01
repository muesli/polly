import Ember from 'ember';
import moment from 'moment';

export function formatTime(params/*, hash*/) {
    return moment(params[0]).tz(moment.tz.guess()).format(params[1]);
}

export default Ember.Helper.helper(formatTime);
