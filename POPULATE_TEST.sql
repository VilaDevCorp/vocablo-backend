CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
DECLARE
    LANG_ID UUID;
    USER_ID UUID;

BEGIN

-- INSERT INTO
--     languages (id, creation_date, code)
-- VALUES
--     (uuid_generate_v4 (), current_timestamp, 'en')
--     RETURNING id INTO LANG_ID;

SELECT id INTO LANG_ID FROM languages WHERE code = 'en';

INSERT INTO
    users (
        id,
        creation_date,
        username,
        email,
        "password",
        validated
    )
VALUES
    (
        uuid_generate_v4 (),
        current_timestamp,
        'johntest',
        'johntest@gmail.com',
        '$2a$14$8nS81fSGzYIMiL3o3u4CUuy7knPG6dPElPALPOVYQysTbADsXeaZS',
        true
    )
    RETURNING id INTO USER_ID;


INSERT INTO user_words
(id, creation_date, term, definitions, learning_progress, language_user_words, user_user_words)
VALUES(
    uuid_generate_v4(),
    current_timestamp - interval '9' minute,
    'water',
    '[{"example": "By the action of electricity, the water was resolved into its two parts, oxygen and hydrogen.", "definition": "A substance (of molecular formula H₂O) found at room temperature and pressure as a clear liquid; it is present naturally as rain, and found in rivers, lakes and seas; its solid form is ice and its gaseous form is steam.", "partOfSpeech": "noun"}, {"example": "The boat was found within the territorial waters.", "definition": "Water in a body; an area of open water.", "partOfSpeech": "noun"}]',
    0,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '8' minute,
    'fire',
    '[{"example": "We sat about the fire singing songs and telling tales.", "definition": "An instance of this chemical reaction, especially when intentionally created and maintained in a specific location to a useful end (such as a campfire or a hearth fire).", "partOfSpeech": "noun"}, {"example": "During hot and dry summers many fires in forests are caused by regardlessly discarded cigarette butts.", "definition": "The occurrence, often accidental, of fire in a certain place, causing damage and danger.", "partOfSpeech": "noun"}]',
    20,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '7' minute,
    'earth',
    '[{"example": "This is good earth for growing potatoes.", "definition": "Soil.", "partOfSpeech": "noun"}, {"example": "She sighed when the plane''s wheels finally touched earth.", "definition": "Any general rock-based material.", "partOfSpeech": "noun"}]',
    40,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '6' minute,
    'wind',
    '[{"example": "As they accelerated onto the motorway, the wind tore the plywood off the car''s roof-rack.", "definition": "Real or perceived movement of atmospheric air usually caused by convection or differences in air pressure.", "partOfSpeech": "noun"}, {"example": "the wind of a cannon ball; the wind of a bellows", "definition": "Air artificially put in motion by any force or action.", "partOfSpeech": "noun"}]',
    60,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '5' minute,
    'pretentious',
    '[{"example": "Her dress was obviously more pretentious than comfortable.", "definition": "Intended to impress others; ostentatious.", "partOfSpeech": "adjective"}, {"example": "Their song titles are pretentious in the context of their basic lyrics.", "definition": "Marked by an unwarranted claim to importance or distinction.", "partOfSpeech": "adjective"}]',
    80,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '4' minute,
    'friend',
    '[{"example": "John and I have been friends ever since we were roommates at college. Trust is important between friends. I used to find it hard to make friends when I was shy.", "definition": "A person other than a family member, spouse or lover whose company one enjoys and towards whom one feels affection.", "partOfSpeech": "noun"}, {"example": "The Automobile Association is every motorist''s friend. The police is every law-abiding citizen''s friend.", "definition": "An associate who provides assistance.", "partOfSpeech": "noun"}, {"example": "a friend of a friend; I added him as a friend on Facebook, but I hardly know him.", "definition": "A person with whom one is vaguely or indirectly acquainted.", "partOfSpeech": "noun"}]',
    100,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '3' minute,
    'laugh',
    '[{"example": "His deep laughs boomed through the room.", "definition": "An expression of mirth particular to the human species; the sound heard in laughing; laughter.", "partOfSpeech": "noun"}, {"definition": "A fun person.", "partOfSpeech": "noun"}]',
    30,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '2' minute,
    'lame',
    '[{"definition": "A stupid or undesirable person.", "partOfSpeech": "noun"}, {"example": "He kept telling these extremely lame jokes all night.", "definition": "Failing to be cool, funny, interesting or relevant.", "partOfSpeech": "adjective"}]',
    50,
    LANG_ID,
    USER_ID
),
(
    uuid_generate_v4(),
    current_timestamp - interval '1' minute,
    'wonderful',
    '[{"example": "They served a wonderful six-course meal.", "definition": "Surprisingly excellent; very good or admirable, extremely impressive.", "partOfSpeech": "adjective"}, {"definition": "Tending to excite wonder; surprising, extraordinary.", "partOfSpeech": "adjective"}]',
    90,
    LANG_ID,
    USER_ID
);

END $$;