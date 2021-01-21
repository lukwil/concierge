CREATE TABLE public.bla (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    name text NOT NULL
);
CREATE TABLE public.distributed_deployment (
    id integer NOT NULL,
    name text NOT NULL,
    name_k8s text,
    container_image text NOT NULL,
    launcher_cpu integer NOT NULL,
    launcher_ram integer NOT NULL,
    worker_count integer NOT NULL,
    worker_cpu integer NOT NULL,
    worker_ram integer NOT NULL,
    worker_gpu integer NOT NULL,
    url_prefix text NOT NULL,
    status_id integer
);
CREATE SEQUENCE public.distributed_deployment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.distributed_deployment_id_seq OWNED BY public.distributed_deployment.id;
CREATE TABLE public.distributed_environment_variables (
    id integer NOT NULL,
    name text NOT NULL,
    value text NOT NULL,
    distributed_deployment_id integer NOT NULL
);
CREATE SEQUENCE public.distributed_environment_variables_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.distributed_environment_variables_id_seq OWNED BY public.distributed_environment_variables.id;
CREATE TABLE public.single_deployment (
    id integer NOT NULL,
    name text NOT NULL,
    name_k8s text,
    cpu integer NOT NULL,
    ram integer NOT NULL,
    gpu integer NOT NULL,
    volume_id integer,
    status_id integer,
    container_image text NOT NULL,
    url_prefix text NOT NULL
);
CREATE SEQUENCE public.single_deployment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.single_deployment_id_seq OWNED BY public.single_deployment.id;
CREATE TABLE public.status (
    id integer NOT NULL,
    name text NOT NULL
);
CREATE SEQUENCE public.status_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.status_id_seq OWNED BY public.status.id;
CREATE TABLE public.tab (
    id integer NOT NULL,
    name text NOT NULL
);
CREATE SEQUENCE public.tab_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.tab_id_seq OWNED BY public.tab.id;
CREATE TABLE public.test (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    name text NOT NULL
);
CREATE TABLE public.volume (
    id integer NOT NULL,
    size integer,
    mount_path text
);
CREATE SEQUENCE public.volume_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.volume_id_seq OWNED BY public.volume.id;
ALTER TABLE ONLY public.distributed_deployment ALTER COLUMN id SET DEFAULT nextval('public.distributed_deployment_id_seq'::regclass);
ALTER TABLE ONLY public.distributed_environment_variables ALTER COLUMN id SET DEFAULT nextval('public.distributed_environment_variables_id_seq'::regclass);
ALTER TABLE ONLY public.single_deployment ALTER COLUMN id SET DEFAULT nextval('public.single_deployment_id_seq'::regclass);
ALTER TABLE ONLY public.status ALTER COLUMN id SET DEFAULT nextval('public.status_id_seq'::regclass);
ALTER TABLE ONLY public.tab ALTER COLUMN id SET DEFAULT nextval('public.tab_id_seq'::regclass);
ALTER TABLE ONLY public.volume ALTER COLUMN id SET DEFAULT nextval('public.volume_id_seq'::regclass);
ALTER TABLE ONLY public.bla
    ADD CONSTRAINT bla_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.distributed_deployment
    ADD CONSTRAINT distributed_deployment_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.distributed_environment_variables
    ADD CONSTRAINT distributed_environment_variables_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.single_deployment
    ADD CONSTRAINT single_deployment_id_key UNIQUE (id);
ALTER TABLE ONLY public.single_deployment
    ADD CONSTRAINT single_deployment_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.status
    ADD CONSTRAINT status_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.tab
    ADD CONSTRAINT tab_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.test
    ADD CONSTRAINT test_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.volume
    ADD CONSTRAINT volume_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.distributed_deployment
    ADD CONSTRAINT distributed_deployment_status_id_fkey FOREIGN KEY (status_id) REFERENCES public.status(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.distributed_environment_variables
    ADD CONSTRAINT distributed_environment_variables_distributed_deployment_id_ FOREIGN KEY (distributed_deployment_id) REFERENCES public.distributed_deployment(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.single_deployment
    ADD CONSTRAINT single_deployment_status_id_fkey FOREIGN KEY (status_id) REFERENCES public.status(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.single_deployment
    ADD CONSTRAINT single_deployment_volume_id_fkey FOREIGN KEY (volume_id) REFERENCES public.volume(id) ON UPDATE RESTRICT ON DELETE CASCADE;
