package hints

const DEFAULT_DICT_NEW = "if '__dict_manager' not in globals():\n    from starkware.cairo.common.dict import DictManager\n    __dict_manager = DictManager()\n\nmemory[ap] = __dict_manager.new_default_dict(segments, ids.default_value)"

const DICT_READ = "dict_tracker = __dict_manager.get_tracker(ids.dict_ptr)\ndict_tracker.current_ptr += ids.DictAccess.SIZE\nids.value = dict_tracker.data[ids.key]"

const DICT_WRITE = "dict_tracker = __dict_manager.get_tracker(ids.dict_ptr)\ndict_tracker.current_ptr += ids.DictAccess.SIZE\nids.dict_ptr.prev_value = dict_tracker.data[ids.key]\ndict_tracker.data[ids.key] = ids.new_value"

const DICT_UPDATE = "# Verify dict pointer and prev value.\ndict_tracker = __dict_manager.get_tracker(ids.dict_ptr)\ncurrent_value = dict_tracker.data[ids.key]\nassert current_value == ids.prev_value, \\\n    f'Wrong previous value in dict. Got {ids.prev_value}, expected {current_value}.'\n\n# Update value.\ndict_tracker.data[ids.key] = ids.new_value\ndict_tracker.current_ptr += ids.DictAccess.SIZE"

const SQUASH_DICT = "dict_access_size = ids.DictAccess.SIZE\naddress = ids.dict_accesses.address_\nassert ids.ptr_diff % dict_access_size == 0, \\\n   'Accesses array size must be divisible by DictAccess.SIZE'\nn_accesses = ids.n_accesses\nif '__squash_dict_max_size' in globals():\n    assert n_accesses <= __squash_dict_max_size, \\\n        f'squash_dict() can only be used with n_accesses<={__squash_dict_max_size}. ' \\n        f'Got: n_accesses={n_accesses}.'\n# A map from key to the list of indices accessing it.\naccess_indices = {}\nfor i in range(n_accesses):\n    key = memory[address + dict_access_size * i]\n    access_indices.setdefault(key, []).append(i)\n# Descending list of keys.\nkeys = sorted(access_indices.keys(), reverse=True)\n# Are the keys used bigger than range_check bound.\nids.big_keys = 1 if keys[0] >= range_check_builtin.bound else 0\nids.first_key = key = keys.pop()"

const SQUASH_DICT_INNER_SKIP_LOOP = "ids.should_skip_loop = 0 if current_access_indices else 1"

const SQUASH_DICT_INNER_FIRST_ITERATION = "current_access_indices = sorted(access_indices[key])[::-1]\ncurrent_access_index = current_access_indices.pop()\nmemory[ids.range_check_ptr] = current_access_index"

const SQUASH_DICT_INNER_CHECK_ACCESS_INDEX = "new_access_index = current_access_indices.pop()\nids.loop_temps.index_delta_minus1 = new_access_index - current_access_index - 1\ncurrent_access_index = new_access_index"

const SQUASH_DICT_INNER_CONTINUE_LOOP = "ids.loop_temps.should_continue = 1 if current_access_indices else 0"

const SQUASH_DICT_INNER_ASSERT_LEN_KEYS = "assert len(keys) == 0"

const SQUASH_DICT_INNER_LEN_ASSERT = "assert len(current_access_indices) == 0"

const SQUASH_DICT_INNER_USED_ACCESSES_ASSERT = "assert ids.n_used_accesses == len(access_indices[key])"

const SQUASH_DICT_INNER_NEXT_KEY = "assert len(keys) > 0, 'No keys left but remaining_accesses > 0.'\nids.next_key = key = keys.pop()"
